/*
 * RUNTIME VULNERABILITY (intentional) — DOM-based cross-site scripting.
 *
 * The visitor's name travels   location.hash  ->  decodeURIComponent  ->  jQuery .html()
 * with no escaping. CodeQL reports this as a DOM XSS data flow, and the site is
 * live on GitHub Pages so the exploit can be shown in the browser:
 *
 *   https://artofthepossible.github.io/ghas-me-up/#<img src=x onerror=alert(document.domain)>
 *
 * Fix (used in the demo's "close the loop" step): $('#greeting').text(...)
 */
(function () {
  'use strict';

  function render() {
    var name = decodeURIComponent(window.location.hash.slice(1)) || 'friend';
    $('#greeting').html('Hello, ' + name + '!');
  }

  $('#hello-form').on('submit', function (event) {
    event.preventDefault();
    window.location.hash = encodeURIComponent($('#name').val());
  });

  window.addEventListener('hashchange', render);
  render();
})();
