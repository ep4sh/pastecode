function Copy(){
  const el = document.createElement('textarea');
  el.value = window.location.href;
  document.body.appendChild(el);
  el.select();
  document.execCommand('copy');
  document.body.removeChild(el);
  alert("Copied link: " + el.value);
}

