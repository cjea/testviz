(function() {
  var filesContainer  = document.getElementById('files-container');
  var files           = document.getElementById('files');
  var nav             = document.getElementById('nav');
  var legend          = document.getElementById('legend');
  var content         = document.getElementById('content');

  var sortDesc = document.getElementById('sortDesc');
  var sortAsc  = document.getElementById('sortAsc');
  var sortName = document.getElementById('sortName');
  var filterEl = document.getElementById('filter');

  if (!files || !content) return;

  var visible;

  function select(part) {
    if (visible) visible.style.display = 'none';
    visible = document.getElementById(part);
    if (!visible) return;
    visible.style.display = 'block';
    location.hash = part;
  }

  files.addEventListener('click', function(e) {
    var t = e.target;
    while (t && t !== files && !(t.classList && t.classList.contains('row'))) { t = t.parentNode; }
    if (!t || t === files) return;

    var part = t.getAttribute('value') || t.getAttribute('data-target');
    if (part) { select(part); window.scrollTo(0,0); }
  }, false);

  function parseCoverage(text) {
    var m = (text || '').match(/([\d.]+)\s*%/);
    return m ? parseFloat(m[1]) : -Infinity;
  }

  function parseName(text) {
    var s = text || '';
    var i = s.lastIndexOf(' (');
    return i >= 0 ? s.slice(0, i) : s;
  }

  function sortRows(cmp) {
    var rows = Array.prototype.slice.call(files.querySelectorAll('.row'));
    rows.sort(cmp);
    rows.forEach(function(r){ files.appendChild(r); });
  }

  if (sortDesc) sortDesc.addEventListener('click', function(){
    sortRows(function(a,b){
      var ac = parseCoverage(a.textContent), bc = parseCoverage(b.textContent);
      if (bc !== ac) return bc - ac;                 // DESC
      var an = parseName(a.textContent), bn = parseName(b.textContent);
      return an.localeCompare(bn);                    // tie by name ASC
    });
  });

  if (sortAsc) sortAsc.addEventListener('click', function(){
    sortRows(function(a,b){
      var ac = parseCoverage(a.textContent), bc = parseCoverage(b.textContent);
      if (ac !== bc) return ac - bc;                  // ASC
      var an = parseName(a.textContent), bn = parseName(b.textContent);
      return an.localeCompare(bn);                    // tie by name ASC
    });
  });

  if (sortName) sortName.addEventListener('click', function(){
    sortRows(function(a,b){
      var an = parseName(a.textContent), bn = parseName(b.textContent);
      var cmp = an.localeCompare(bn);                 // ASC
      if (cmp) return cmp;
      var ac = parseCoverage(a.textContent), bc = parseCoverage(b.textContent);
      return bc - ac;                                 // tie by coverage DESC
    });
  });

 function applyFilter() {
    var q = (filterEl && filterEl.value != null) ? filterEl.value : "";
    var re = null;
    if (q.length) {
      try { re = new RegExp(q); } catch (e) { re = null; } // invalid regex => show all
    }
    var rows = files.querySelectorAll('.row');
    for (var i = 0; i < rows.length; i++) {
      var txt = rows[i].textContent || "";
      var show = !re || re.test(txt);
      rows[i].style.display = show ? "" : "none";
    }
  }
  if (filterEl) {
    filterEl.addEventListener('input', applyFilter);
  }

  function applyMode() {
    function show(el, val) { if (el) el.style.display = val }
    var hasFile = !!(location.hash && location.hash.length > 1);

    if (hasFile) {
      var part = location.hash.substr(1);
      if (!visible || visible.id !== part) select(part);
      show(filesContainer, 'none')
      show(legend, 'block')
      show(content, 'block')
    } else {
      if (visible) { visible.style.display = 'none'; visible = null; }
      show(visible, 'none'); visible = null;
      show(filesContainer, 'block');
      show(legend, 'none');
      show(content, 'none');
    }
  }

  window.addEventListener('hashchange', applyMode);
  applyMode();
  applyFilter();
})();
