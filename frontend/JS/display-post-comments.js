import {main} from './layout.js'
import {handleAddComment} from './comments-service.js'
import {loadPosts, buildMain} from './layout.js'

let currentPost = {}
let comments = {}
let author = {}

async function retrievePost(postId) {
  try {
    const res = await fetch(`http://localhost:8080/post?id=${postId}`)
    if (!res.ok) {
      console.error('Erreur requête retrievePost :', res.status)
    }

    const data = await res.json()
    currentPost = data.post
    comments = data.comments
    author = data.author

    return (currentPost, comments, author)
  } catch (err) {
    console.error('Pas de réponse de PostHandler :', err)
  }
}

async function postLayout(postId) {
  await retrievePost(postId)
  main.innerHTML = `
    <div id="post-card" class="card">
	<span>Créateur: <strong>${author.username}</strong></span>
  <span>${new Date(currentPost.createdat).toLocaleString('fr-FR')}</span>
      <span>Titre: <strong>${currentPost.title}</strong></span>
	  <span>${currentPost.content}</span>
    </div>		
		<div id="comments-card" class="card"><span>Commentaires:</span></div>
	<form id="comment-form">
		<textarea id="comment-text" name="comment" placeholder="Ajouter un commentaire" required maxlength="500"></textarea>
		<button id="comment-btn" type="submit">Publier</button>
		<p id="error"></p>
	</form>
	<button id="returnHome">Retour à l'accueil</div>
    `

  document.getElementById('returnHome').addEventListener('click', async () => {
    const posts = await loadPosts()
    buildMain(posts)
  })

  const commentsSection = document.getElementById('comments-card')

  if (comments) {
    comments.forEach((comment) => {
      const div = document.createElement('div')
      div.classList.add('comment-row')

      div.innerHTML = `
    <span>${new Date(comment.createdat).toLocaleString('fr-FR')}</span>
		<span><strong>${comment.authorname}:</strong></span>
		<span>${comment.content}</span>
`
      commentsSection.appendChild(div)
    })
  } else {
    const div = document.createElement('div')
    div.classList.add('comment-row')

    div.innerHTML = `
			<span>Aucun commentaire publié pour ce post</span>
		`
    commentsSection.appendChild(div)
  }

  document
    .getElementById('comment-form')
    .addEventListener('submit', handleAddComment)
}

export {postLayout, currentPost, comments, author}
