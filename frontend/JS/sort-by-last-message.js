export async function sortLastMessage() {
  try {
    const res = await fetch('/messages/lastUser')
    if (!res.ok) throw new Error(`Erreur: ${res.status}`)
    const data = await res.json()
    console.log('retour de fetch lastuser', data)
  } catch (err) {
    console.error('Erreur fetch:', err)
  }
}
