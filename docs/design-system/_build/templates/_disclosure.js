<script>
/* Open in place (v1.4): pillar index and compact tickets. No dependency. */
(function () {
  /* Pillars: one open at a time */
  var heads = document.querySelectorAll('.pillar-index__head');
  function parts(h) { var b = h.querySelector('.pillar-index__btn'); return [b, document.getElementById(b.getAttribute('aria-controls'))]; }
  heads.forEach(function (h) {
    var btn = parts(h)[0];
    btn.addEventListener('click', function () {
      var open = btn.getAttribute('aria-expanded') !== 'true';
      heads.forEach(function (o) {
        var p = parts(o), on = o === h && open;
        p[0].setAttribute('aria-expanded', String(on)); o.classList.toggle('is-open', on); p[1].classList.toggle('is-open', on);
      });
    });
  });
  /* Tickets: from 900px up the tickets of one row open together, so the row stays level and the scopes compare */
  var wide = window.matchMedia('(min-width: 900px)');
  function setTicket(t, open) {
    t.setAttribute('aria-expanded', String(open));
    document.getElementById(t.getAttribute('aria-controls')).classList.toggle('is-open', open);
    t.closest('.ticket').classList.toggle('is-open', open);
  }
  document.querySelectorAll('.ticket__toggle').forEach(function (t) {
    t.addEventListener('click', function () {
      var open = t.getAttribute('aria-expanded') !== 'true';
      var row = wide.matches ? t.closest('.grid').querySelectorAll('.ticket__toggle') : [t];
      row.forEach(function (o) { setTicket(o, open); });
    });
  });
})();
</script>
