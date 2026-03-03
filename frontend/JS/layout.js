import {renderCreatePost, loadCategories} from './post-service.js'
import {Logout} from './authentication.js'
import {postLayout} from './display-post-comments.js'

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
      console.log("Réponse brute :", data);
      console.log('Posts reçus:', data.allposts)
      buildMain(data.allposts)
    } catch (err) {
      console.error('Erreur fetch posts:', err)
    }
  })

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


// ✅ Charger tous les utilisateurs depuis l'API
async function loadAllUsers() {
  const usersList = document.querySelector('.users-list')

  if (!usersList) {
    console.error('❌ .users-list introuvable dans le DOM !')
    return
  }

  console.log('🔄 Chargement des utilisateurs...')

  try {
    const response = await fetch('/api/users')
    if (!response.ok) throw new Error('Erreur récupération utilisateurs')

    const allUsers = await response.json()
    console.log('✅ Utilisateurs chargés:', allUsers)

    usersList.innerHTML = ''

    if (allUsers.length === 0) {
      console.warn('⚠️ Aucun utilisateur trouvé dans la base')
      usersList.innerHTML = '<p>Aucun utilisateur</p>'
      return
    }

    allUsers.forEach((user) => {
      const userEl = document.createElement('div')
      userEl.classList.add('user-item')

      // ✅ Par défaut, tous sont hors ligne (classe .offline)
      userEl.classList.add('offline')

      userEl.textContent = user.nickname
      userEl.dataset.userId = user.id

      userEl.addEventListener('click', () => {
        userEl.classList.remove('has-notification')
        handleChatClick(null, user.id, user.nickname)
      })

      usersList.appendChild(userEl)
    })

    console.log(
      `✅ ${allUsers.length} utilisateurs affichés (hors ligne par défaut)`
    )
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

  // ✅ Créer un Set des IDs en ligne
  const onlineUserIds = new Set(onlineUsers.map((u) => u.id))
  console.log('🟢 Utilisateurs en ligne:', Array.from(onlineUserIds))

  // ✅ Parcourir tous les .user-item et mettre à jour leur statut
  const userItems = usersList.querySelectorAll('.user-item')

  if (userItems.length === 0) {
    console.warn('⚠️ Aucun .user-item trouvé pour mise à jour')
    return
  }

  userItems.forEach((userEl) => {
    const userId = parseInt(userEl.dataset.userId)

    if (onlineUserIds.has(userId)) {
      // 🟢 En ligne → retirer .offline
      userEl.classList.remove('offline')
      console.log(`🟢 ${userEl.textContent} est EN LIGNE`)
    } else {
      // 🔴 Hors ligne → ajouter .offline
      userEl.classList.add('offline')
      console.log(`🔴 ${userEl.textContent} est HORS LIGNE`)
    }
  })
}

async function buildMain(posts = []) {
  if (!Array.isArray(posts)) posts = [];
  const res = await fetch("/categories")
  const categories = await res.json();
  console.log("categories:", categories);

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

    const categoryNames = (post.category_ids || [])
    .map(id => catMap[id])
    .join(', ');

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
  console.log(data.allposts)
  return data.allposts
}

async function showApp() {
  document.getElementById('auth-container').style.display = 'none'
  document.getElementById('app-container').style.display = 'contents'
  buildHeader()

  const posts = await loadPosts()
  console.log(posts)
  buildMain(posts)
}

export { header, main, sideBar, buildHeader, showApp, loadPosts, buildMain };
