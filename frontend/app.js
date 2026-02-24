import {ToggleRegisterFrom} from "./JS/connection.js";
import {CreateAccount} from "./JS/create-account.js";
import {Authentication} from "./JS/authentication.js";
import {showApp} from "./JS/layout.js";
import { initChat } from "./JS/chat.js";

const authContainer = document.getElementById("auth-container");
const appContainer = document.getElementById("app-container");

function showLogin() {
  authContainer.style.display = "flex";
  appContainer.style.display = "none";
}

// Quand le HTML est entièrement chargé, exécute cette fonction
// Garantit que le JS s'exécute une fois que la page est prête
document.addEventListener("DOMContentLoaded", () => {
  // Préparation des formulaires
  showLogin();
  ToggleRegisterFrom();
  CreateAccount();
  Authentication();

  // On vérifie si une session existe déjà, envoie de la requête HTTP GET
  fetch("/auth/me", {
    credentials: "include", // Dit au navigateur d'envoyer aussi les cookies
  })
    .then((res) => {
      if (!res.ok) throw new Error("Not logged in");
      return res.json();
    })
    .then((user) => {
      showApp();
      document.getElementById("welcome-message").textContent =
        `Bienvenue, ${user.nickname} !`;
      initChat("miaoutest");
    })
    .catch(() => {
      // Pas connecté : on reste sur le login, déjà initialisé
    });
});
