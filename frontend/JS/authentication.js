import { showApp } from "./layout.js";

export function Authentication() {
	const connectionForm = document.querySelector(".connection-form");
	console.log("form récup");

	connectionForm.addEventListener("submit", (e) => {
		e.preventDefault(); // empêche le rechargement de la page

		const login = document.getElementById("register-username").value;
		const password = document.getElementById("password").value;

		console.log("username et password recup");

		fetch("/auth/login", {
			// envoie la requête au serveur sur la route /auth/login
			method: "POST", // on envoie des données
			headers: {
				"Content-Type": "application/json", // le serveur sait qu’on envoie du JSON
			},
			credentials: "include",
			body: JSON.stringify({
				// on transforme nos valeurs JS en JSON
				login: login,
				password: password,
			}),
		})
			.then((response) => response.json()) // le serveur renvoie du JSON, on le transforme en objet JS
			.then((data) => {
				// Gestion de la réponse (succès ou erreur)
				console.log("Réponse serveur :", data);
				if (data.success) {
					showApp();

					const registerForm = document.querySelector(".register-form");
					const forumSection = document.querySelector(".forum-section");
					const welcomeMessage = document.getElementById("welcome-message");

					// Masquer les formulaires
					// Afficher la section forum
					connectionForm.style.display = "none";
					registerForm.style.display = "none";
					forumSection.style.display = "flex";

					// Mettre le pseudo dans le message
					welcomeMessage.textContent = `Bienvenue, ${data.user.nickname} !`;
				} else {
					alert("Utilisateur inconnu ou mauvais mot de passe");
				}
			})
			.catch((error) => {
				// attrape les erreurs réseau
				console.log("Erreur réseau :", error);
			});

		console.log("Login :", login, "Password :", password);
	});
}

export function Logout() {
	const forumSection = document.querySelector(".forum-section");
	const authContainer = document.getElementById("auth-container");
	const appContainer = document.getElementById("app-container");
	const connectionForm = document.querySelector(".connection-form");
	const registerForm = document.querySelector(".register-form");

	if (forumSection) forumSection.style.display = "none";
	appContainer.style.display = "none";
	authContainer.style.display = "flex";
	connectionForm.style.display = "flex";
	registerForm.style.display = "none";

	document.getElementById("register-username").value = "";
	document.getElementById("password").value = "";
}
