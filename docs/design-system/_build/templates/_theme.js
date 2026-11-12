<script>
(function () {
  var root = document.documentElement;
  function get() { try { return localStorage.getItem('ivm-theme'); } catch (e) { return null; } }
  function set(t) { try { localStorage.setItem('ivm-theme', t); } catch (e) {} }
  var t = get(); if (t) root.dataset.theme = t;
  document.addEventListener('click', function (e) {
    var b = e.target.closest('[data-toggle-theme]'); if (!b) return;
    var dark = root.dataset.theme ? root.dataset.theme === 'dark' : matchMedia('(prefers-color-scheme: dark)').matches;
    root.dataset.theme = dark ? 'light' : 'dark'; set(root.dataset.theme);
  });
})();
</script>
