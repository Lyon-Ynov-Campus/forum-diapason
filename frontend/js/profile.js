function timeAgo(dateStr) {
    const diff = Math.floor((Date.now() - new Date(dateStr)) / 60000)
    if (diff < 60) return `il y a ${diff} minute${diff > 1 ? 's' : ''}`
    const h = Math.floor(diff / 60)
    if (h < 24) return `il y a ${h} heure${h > 1 ? 's' : ''}`
    const d = Math.floor(h / 24)
    return `il y a ${d} jour${d > 1 ? 's' : ''}`
}

function createPostCard(post) {
    const card = document.getElementById('post-card').content.cloneNode(true)
    const authorEl = card.querySelector('.post-author')
    authorEl.textContent = post.author_pseudo
    authorEl.href = `/profile?id=${post.user_id}`
    card.querySelector('.post-title').textContent = post.titre
    card.querySelector('.post-content').textContent = post.contenu
    card.querySelector('.post-date').textContent = timeAgo(post.date_publication)
    card.querySelector('.post-tags').textContent = (post.tags || []).map(t => `#${t}`).join(' ')
    card.querySelector('.post-likes').textContent = post.like_count

    if (post.image_url) {
        const img = document.createElement('img')
        img.src = `http://localhost:8080${post.image_url}`
        img.className = 'max-w-full max-h-96 object-contain object-left mt-1'
        img.alt = post.titre
        const titleEl = card.querySelector('.post-title')
        titleEl?.insertAdjacentElement('afterend', img)
    }

    const article = card.querySelector('article')
    article.style.cursor = 'pointer'
    article.addEventListener('click', (e) => {
        if (e.target.closest('button')) return
        window.location.href = `/post?id=${post.id}`
    })
    initPostCard(article, post)
    return card
}

const params      = new URLSearchParams(window.location.search)
const profileId   = parseInt(params.get('id')) || null
const profilePseudo = params.get('pseudo') || null

const userId = profileId || 1

const loadProfilePosts = () => fetch(`${API}/api/users/${userId}/posts`, { credentials: 'include' })
    .then(r => r.json())
    .then(posts => {
    const container = document.getElementById('profile-posts-container')
    if (!container) return

    const noPostsEl = document.getElementById('profile-no-posts')

    if (!posts || posts.length === 0) return

    // Cache le message "aucun post"
    if (noPostsEl) noPostsEl.remove()

    ;(posts || []).forEach(post => container.appendChild(createPostCard(post)))

    document.getElementById('open-edit-profile-btn')?.addEventListener('click', () => {
        openEditProfileModal(currentUser || {})
    })
})

if (typeof checkAuth === 'function') checkAuth().then(loadProfilePosts)
else loadProfilePosts()
