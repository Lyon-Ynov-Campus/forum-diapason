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
        img.className = 'w-full object-cover max-h-48 mt-1'
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

function createTopPostCard(post) {
    const card = document.getElementById('top-post-card').content.cloneNode(true)
    const a = card.querySelector('.top-post-author')
    a.textContent = post.author_pseudo
    a.href = `/profile?id=${post.user_id}`
    card.querySelector('.top-post-content').textContent = post.contenu
    card.querySelector('.top-post-date').textContent = timeAgo(post.date_publication)
    return card
}

function loadPosts() {
    fetch(`${API}/api/posts`, { credentials: 'include' })
        .then(r => r.json())
        .then(posts => {
            if (!Array.isArray(posts)) return
            const container = document.getElementById('posts-container')
            posts.forEach(post => container?.appendChild(createPostCard(post)))
        })

    fetch(`${API}/api/posts/top?limit=6`, { credentials: 'include' })
        .then(r => r.json())
        .then(posts => {
            if (!Array.isArray(posts)) {
                console.error('top posts error:', posts)
                return
            }
            const topContainer = document.getElementById('top-posts-container')
            posts.forEach(post => topContainer?.appendChild(createTopPostCard(post)))
        })
        .catch(err => console.error('top posts fetch error:', err))
}

function reloadPosts() {
    const container    = document.getElementById('posts-container')
    const topContainer = document.getElementById('top-posts-container')
    if (container)    container.innerHTML    = ''
    if (topContainer) topContainer.innerHTML = ''
    loadPosts()
}

loadPosts()
