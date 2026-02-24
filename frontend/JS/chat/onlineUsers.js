// handleOnlineUsers Gérer la liste des utilisateurs en ligne
export function handleOnlineUsers(users, main, ws) {
  // On récupère le conteneur de la liste des utilisateurs
  const usersList = document.querySelector('.users-list')

  // Sécurité : si l'élément n'existe pas, on arrête
  if (!usersList) {
    console.warn('.users-list introuvable')
    return
  }

  // On vide la liste actuelle (évite les doublons)
  usersList.innerHTML = ''

  // Pour chaque utilisateur reçu depuis le serveur
  users.forEach((user) => {
    // Création d’un élément HTML pour l’utilisateur
    const userEl = document.createElement('div')

    // Ajout d’une classe CSS
    userEl.classList.add('user-item')

    // Affichage du nom dans l’élément
    userEl.textContent = user.name

    // Stockage de l’id utilisateur dans un data-attribute
    userEl.dataset.userId = user.id

    // Événement : clic sur un utilisateur
    userEl.addEventListener('click', () => {
      // Supprime l’indicateur de notification si présent
      userEl.classList.remove('has-notification')

      // Met à jour l’utilisateur actuellement sélectionné
      currentChatUserId = user.id

      // Ouvre la conversation avec cet utilisateur
      openConversation(main, ws, user.id, user.name)
    })

    // Ajoute l’utilisateur dans la liste affichée
    usersList.appendChild(userEl)
  })

  // Log informatif
  console.log(`${users.length} utilisateurs affichés`)
}
