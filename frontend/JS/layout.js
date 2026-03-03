import {renderCreatePost, loadCategories} from './post-service.js'
import {Logout} from './authentication.js'
import {postLayout} from './display-post-comments.js'
import {openChatWith} from './chat.js'

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
    <select id="category" name="category" required>
          <option value="all">Tous les posts</option>
        </select>
  </nav>
  <div class="forum-section">
    <div class="profile-section">
      <p id="welcome-message"></p>
      <img src="./frontend/img/profil.gif" alt="image profil" class="profil-icon">
      </div>
      <button id="logoutBtn">Déconnexion</button>
    </div>
`
  loadCategories()

  const categorySelect = document.getElementById("category");
  categorySelect.addEventListener("change", async () => {
    const category = categorySelect.value;
    try {
      const res = await fetch(`/post?id=0&category=${category}`)
      if (!res.ok) {
        console.error('Erreur fetch posts:', res.status, await res.text())
        return
      }
      const data = await res.json()
      buildMain(data.allposts)
    } catch (err) {
      console.error('Erreur fetch posts:', err)
    }
  })

  document.getElementById('logoutBtn').addEventListener('click', async () => {
    try {
      const response = await fetch('/logout', { method: 'POST', credentials: 'include' })
      if (response.ok) {
        alert('Déconnecté avec succès !')
      } else {
        alert('Erreur lors de la déconnexion.')
      }
    } catch (err) {
      console.error('Erreur fetch logout :', err)
    }
  })

  document.getElementById('new-post-btn').addEventListener('click', renderCreatePost)
  document.getElementById('logoutBtn').addEventListener('click', Logout)
  document.getElementById('home-btn').addEventListener('click', async () => {
    const posts = await loadPosts()
    buildMain(posts)
  })
}

// ✅ Construire la sidebar
async function buildSidebar() {
  sideBar.innerHTML = `<h2>Utilisateurs</h2>
  <div class="users-list"></div>`

  await loadAllUsers() // ✅ attendre que les users soient dans le DOM

  // ✅ Appel immédiat
  fetch("/online-users")
    .then(res => res.json())
    .then(updateUsersList)

  // ✅ Polling toutes les 3s APRÈS que .users-list est peuplé
  setInterval(() => {
    fetch("/online-users")
      .then(res => res.json())
      .then(updateUsersList)
  }, 3000)
}

// ✅ Charger tous les utilisateurs depuis l'API
async function loadAllUsers() {
  const usersList = document.querySelector('.users-list')
  if (!usersList) {
    console.error('❌ .users-list introuvable dans le DOM !')
    return
  }

  try {
    const response = await fetch('/api/users')
    if (!response.ok) throw new Error('Erreur récupération utilisateurs')

    const allUsers = await response.json()
    usersList.innerHTML = ''

    if (allUsers.length === 0) {
      usersList.innerHTML = '<p>Aucun utilisateur</p>'
      return
    }

    allUsers.forEach((user) => {
      const userEl = document.createElement('div')
      userEl.classList.add('user-item', 'offline')
      userEl.dataset.userId = user.id
      userEl.dataset.userName = user.nickname // ✅ pour la comparaison

      // ✅ Point de statut
      const dot = document.createElement('span')
      dot.classList.add('status-dot', 'dot-red')
      userEl.appendChild(dot)

      const name = document.createElement('span')
      name.textContent = user.nickname
      userEl.appendChild(name)

      userEl.addEventListener('click', () => {
        userEl.classList.remove('has-notification')
        openChatWith(user.nickname)
      })

      usersList.appendChild(userEl)
    })
  } catch (error) {
    console.error('❌ Erreur chargement utilisateurs:', error)
  }
}

// ✅ Mettre à jour les statuts (en ligne/hors ligne)
function updateUsersList(onlineUsers) {
  const usersList = document.querySelector('.users-list')
  if (!usersList) {
    console.error('❌ .users-list introuvable pour mise à jour')
    return
  }

  const onlineUserIds = new Set(onlineUsers.map((u) => u.id))

  const userItems = usersList.querySelectorAll('.user-item')
  userItems.forEach((userEl) => {
    const userId = parseInt(userEl.dataset.userId)
    const dot = userEl.querySelector('.status-dot')

    if (onlineUserIds.has(userId)) {
      userEl.classList.remove('offline')
      userEl.classList.add('online')
      if (dot) { dot.classList.remove('dot-red'); dot.classList.add('dot-green') }
    } else {
      userEl.classList.remove('online')
      userEl.classList.add('offline')
      if (dot) { dot.classList.remove('dot-green'); dot.classList.add('dot-red') }
    }
  })
}

async function buildMain(posts = []) {
  if (!Array.isArray(posts)) posts = [];
  const res = await fetch("/categories")
  const categories = await res.json();

  const catMap = {};
  categories.forEach(c => catMap[c.id] = c.name);

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

    const categoryNames = (post.category_ids || []).map(id => catMap[id]).join(', ');
    div.innerHTML = `
      <span>${truncate(post.title, 30)}</span>
      <span>${categoryNames}</span>
      <span>${truncate(post.content, 50)}</span>
    `
    div.addEventListener('click', () => postLayout(post.id))
    list.appendChild(div)
  })
}

function truncate(text, max = 50) {
  if (!text) return "";
  return text.length > max ? text.slice(0, max) + "..." : text;
}

async function loadPosts() {
  const res = await fetch('/post?id=0')
  if (!res.ok) {
    console.error('Erreur lors du chargement des posts')
    return []
  }
  const data = await res.json()
  return data.allposts
}

async function showApp() {
  document.getElementById('auth-container').style.display = 'none'
  document.getElementById('app-container').style.display = 'contents'
  buildHeader()
  await buildSidebar()

  const posts = await loadPosts()
  buildMain(posts)
}

export { header, main, sideBar, buildHeader, showApp, loadPosts, buildMain };