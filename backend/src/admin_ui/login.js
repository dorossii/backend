(function () {
  "use strict";

  var form = document.getElementById("login-form");
  var errorEl = document.getElementById("login-error");

  form.addEventListener("submit", function (event) {
    event.preventDefault();
    errorEl.hidden = true;
    errorEl.textContent = "";

    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;

    fetch("/admin/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "same-origin",
      body: JSON.stringify({ Username: username, Password: password }),
    })
      .then(function (res) {
        if (!res.ok) {
          return res.json().then(function (body) {
            throw new Error((body && body.error) || "ログインに失敗しました");
          });
        }
        window.location.href = "/admin/";
      })
      .catch(function (err) {
        errorEl.textContent = err.message;
        errorEl.hidden = false;
      });
  });
})();
