import {postLayout} from './display-post-comments.js'
import {renderCreatePost} from './post-service.js'

const header = document.getElementById('header')
const main = document.getElementById('main-content')

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

  document.getElementById('home-btn').addEventListener('click', async () => {
    const posts = await loadPosts()
    buildMain(posts)
  })
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

  const posts = await loadPosts()
  console.log(posts)
  buildMain(posts)
}

export {header, main, buildHeader, showApp}
