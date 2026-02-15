// Add church code to HTMX requests
document.addEventListener('htmx:configRequest', (event) => {
  const code = localStorage.getItem('last_church');
  console.log(code);
  event.detail.headers['X-Church-Code'] = code;
});



