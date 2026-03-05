import {main, showApp, loadPosts, buildMain} from './layout.js'

let isUsed = false

function renderCreatePost() {
  console.log('affichage formulaire post')

  main.innerHTML = `
    <h2>Nouveau post</h2>
    <form id="post-form">
      <div class="title">
        <label for="title">Titre</label>
        <textarea id="title" name="title" required maxlength="30"></textarea>
      </div>

      <div class="message">
        <label for="message">Message</label>
        <textarea id="message" name="message" placeholder="Text" required maxlength="500"></textarea>
      </div>

      <div class="select-category">
        <label for="category">Catégorie</label>
        <select id="post-category" name="category" required>
          <option value="" disabled selected>Choisir une catégorie</option>
        </select>
      </div>

      <div class="newPost">
        <button id="newPost" type="submit">Publier</button>
      </div>

      <p id="error"></p>
    </form>
    <button id="returnHome">Retour à l'accueil</div>
  `

  loadCategories()
  document.getElementById('returnHome').addEventListener('click', async () => {
    const posts = await loadPosts()
    buildMain(posts)
  })
  document
    .getElementById('post-form')
    .addEventListener('submit', handleCreatePost)
}

export async function loadCategories() {
  try {
    const res = await fetch('/categories')
    const categories = await res.json()

    const headerSelect = document.getElementById('category')
    const formSelect = document.getElementById('post-category')

    // Remplir le select du header
    if (headerSelect) {
      headerSelect.innerHTML = `<option value="all">Tous les posts</option>`
      categories.forEach((cat) => {
        const option = document.createElement('option')
        option.value = cat.id
        option.textContent = cat.name
        headerSelect.appendChild(option)
      })
    }

    // Remplir le select du formulaire
    if (formSelect) {
      formSelect.innerHTML = `<option value="" disabled selected>Choisir une catégorie</option>`
      categories.forEach((cat) => {
        const option = document.createElement('option')
        option.value = cat.id
        option.textContent = cat.name
        formSelect.appendChild(option)
      })
    }
  } catch (err) {
    console.error('Erreur chargement catégories :', err)
  }
}

function getSelectedCategories(form) {
  const categorySelect = document.getElementById('post-category')

  if (!categorySelect.value) return []

  return [Number(categorySelect.value)]
}

async function handleCreatePost(e) {
  isUsed = false
  e.preventDefault()
  console.log('handleCreatePost')

  const form = e.target
  let categoriesId = getSelectedCategories(form)

  const data = {
    title: form.title.value,
    content: form.message.value,
    category_ids: categoriesId
  }

  console.log(data)

  const res = await fetch('/post', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(data)
  })

  const result = await res.json()

  if (!res.ok) {
    document.getElementById('error').textContent = result.error
  } else {
    alert('Post créé !')
    showApp()
  }
}

export {renderCreatePost}
