import {
	postLayout,
	currentPost,
	comments,
	author,
} from "./display-post-comments.js";

async function handleAddComment(e) {
	e.preventDefault();
	const form = e.target;

	const data = {
		postid: currentPost.id,
		authorid: author.id,
		content: form.comment.value,
	};

	const res = await fetch("/comment", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(data),
	});

	const result = await res.json();

	if (!res.ok) {
		document.getElementById("error").textContent = result.error;
	} else {
		alert("Commentaire publié !");
		await postLayout(currentPost.id);
	}
}

export { handleAddComment };
