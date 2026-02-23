// Gérer la liste des utilisateurs en ligne
export function handleOnlineUsers(users, main, ws) {
  const usersList = document.querySelector('.users-list')
  if (!usersList) {
    console.warn('.users-list introuvable')
    return
  }

  usersList.innerHTML = ''

  users.forEach((user) => {
    const userEl = document.createElement('div')
    userEl.classList.add('user-item')
    userEl.textContent = user.name
    userEl.dataset.userId = user.id

    userEl.addEventListener('click', () => {
      userEl.classList.remove('has-notification')
      currentChatUserId = user.id
      openConversation(main, ws, user.id, user.name)
    })

    usersList.appendChild(userEl)
  })

  console.log(`${users.length} utilisateurs affichés`)
}
