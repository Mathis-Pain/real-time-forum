import {Logout} from './authentication.js'
import {handleChatClick} from './chat/chat.js'
import {initWebSocket, addMessageHandler} from './chat/websocket.js'
import {postLayout} from './display-post-comments.js'
import {renderCreatePost} from './post-service.js'

const header = document.getElementById('header')
const main = document.getElementById('main-content')
const sideBar = document.getElementById('sidebar')

function buildHeader() {
  header.innerHTML = `<div class="header-left">
    <h1>Real time forum</h1>
  </div>
  <nav class="header-nav">
    <button id="new-post-btn">Nouveau post</button>
    <button id="home-btn">Home</button>
    <button id="categories-btn">Catégories</button>
  </nav>
  <div class="forum-section">
    <div class="profile-section">
      <p id="welcome-message"></p>
      <img src="./frontend/img/profil.gif" alt="image profil" class="profil-icon">
    </div>
    <button id="logoutBtn">Déconnexion</button>
  </div>`

  document.getElementById('logoutBtn').addEventListener('click', async () => {
    try {
      const response = await fetch('/logout', {
        method: 'POST',
        credentials: 'include'
      })

      if (response.ok) {
        alert('Déconnecté avec succès !')
      } else {
        alert('Erreur lors de la déconnexion.')
      }
    } catch (err) {
      console.error('Erreur fetch logout :', err)
    }
  })

  const postBtn = document.getElementById('new-post-btn')
  postBtn.addEventListener('click', renderCreatePost)

  const logoutBtn = document.getElementById('logoutBtn')
  logoutBtn.addEventListener('click', Logout)

  document.getElementById('home-btn').addEventListener('click', async () => {
    const posts = await loadPosts()
    buildMain(posts)
  })
}

export function buildSidebar() {
  sideBar.innerHTML = `<h2>Utilisateurs</h2>
  <div class="users-list"></div>`

  const ws = initWebSocket()

  ws.addEventListener('open', () => {
    console.log('WebSocket prêt, chargement utilisateurs...')
    loadAllUsers()
  })

  if (ws.readyState === WebSocket.OPEN) {
    loadAllUsers()
  }

  addMessageHandler(handleSidebarMessages)
}

function handleSidebarMessages(data) {
  if (data.type === 'message') {
    const userItem = document.querySelector(
      `.user-item[data-user-id="${data.sender_id}"]`
    )
    if (userItem) {
      userItem.classList.add('has-notification')
    }
  }

  if (data.type === 'online_users') {
    console.log('Mise à jour utilisateurs en ligne:', data.users)
    updateSidebarOnlineStatus(data.users)
  }
}

function updateSidebarOnlineStatus(onlineUsers) {
  const onlineIds = new Set(onlineUsers.map((u) => u.id))

  const userItems = document.querySelectorAll('#sidebar .user-item')

  if (userItems.length === 0) {
    console.warn('Aucun user-item dans la sidebar')
    return
  }

  userItems.forEach((el) => {
    const userId = parseInt(el.dataset.userId)

    if (onlineIds.has(userId)) {
      el.classList.remove('offline')
      el.classList.add('online')
    } else {
      el.classList.remove('online')
      el.classList.add('offline')
    }
  })

  console.log(
    `${onlineIds.size} utilisateurs en ligne sur ${userItems.length} total`
  )
}

async function loadAllUsers() {
  const usersList = document.querySelector('.users-list')

  if (!usersList) {
    console.error('.users-list introuvable dans le DOM !')
    return
  }

  console.log('Chargement des utilisateurs...')

  try {
    const response = await fetch('/api/users')
    if (!response.ok) throw new Error('Erreur récupération utilisateurs')

    const allUsers = await response.json()
    console.log('Utilisateurs chargés:', allUsers)

    usersList.innerHTML = ''

    if (allUsers.length === 0) {
      console.warn('Aucun utilisateur trouvé dans la base')
      usersList.innerHTML = '<p>Aucun utilisateur</p>'
      return
    }

    allUsers.forEach((user) => {
      console.log(`Utilisateur: ${user.nickname} (online: ${user.online})`)
      const userEl = document.createElement('div')
      userEl.classList.add('user-item')
      userEl.textContent = user.nickname
      userEl.dataset.userId = user.id

      userEl.addEventListener('click', () => {
        userEl.classList.remove('has-notification')
        handleChatClick(user.id, user.nickname)
      })

      usersList.appendChild(userEl)
    })

    console.log(`${allUsers.length} utilisateurs affichés`)
  } catch (error) {
    console.error('Erreur chargement utilisateurs:', error)
  }
}

function buildMain(posts = []) {
  main.innerHTML = `
    <h2>Posts</h2>
    <div class="posts-header">
      <span>Titre</span>
      <span>Catégorie(s)</span>
      <span>Texte</span>
    </div>
    <div id="posts-list"></div>
  `

  const list = document.getElementById('posts-list')

  posts.forEach((post) => {
    const div = document.createElement('div')
    div.classList.add('posts-row')

    div.innerHTML = `
      <span>${post.title}</span>
      <span>${(post.category_ids || []).join(', ')}</span>
      <span>${post.content}</span>
    `
    div.addEventListener('click', () => postLayout(post.id))

    list.appendChild(div)
  })
}

async function loadPosts() {
  const res = await fetch('/post?id=0')
  if (!res.ok) {
    console.error('Erreur lors du chargement des posts')
    return []
  }
  const data = await res.json()
  console.log(data.allposts)
  return data.allposts
}

async function showApp() {
  document.getElementById('auth-container').style.display = 'none'
  document.getElementById('app-container').style.display = 'contents'
  buildHeader()
  buildSidebar()

  const posts = await loadPosts()
  console.log(posts)
  buildMain(posts)
}

export {header, main, sideBar, buildHeader, showApp}
