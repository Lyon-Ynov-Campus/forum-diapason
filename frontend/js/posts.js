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
    card.querySelector('.post-author').textContent = post.author
    card.querySelector('.post-title').textContent = post.title
    card.querySelector('.post-content').textContent = post.content
    card.querySelector('.post-date').textContent = timeAgo(post.created_at)
    card.querySelector('.post-tags').textContent = (post.tags || []).map(t => `#${t}`).join(' ')
    card.querySelector('.post-likes').textContent = post.likes

    const article = card.querySelector('article')
    article.style.cursor = 'pointer'
    article.addEventListener('click', (e) => {
        if (e.target.closest('button')) return 
        window.location.href = `/post?id=${post.id}`
    })

    return card
}

function createTopPostCard(post) {
    const card = document.getElementById('top-post-card').content.cloneNode(true)
    card.querySelector('.top-post-author').textContent = post.author
    card.querySelector('.top-post-content').textContent = post.content
    card.querySelector('.top-post-date').textContent = timeAgo(post.created_at)
    return card
}

fetch('/data/posts.json')
    .then(r => r.json())
    .then(posts => {
        const container = document.getElementById('posts-container')
        posts.forEach(post => container.appendChild(createPostCard(post)))

        const topContainer = document.getElementById('top-posts-container')
        ;[...posts]
            .sort((a, b) => b.likes - a.likes)
            .slice(0, 6)
            .forEach(post => topContainer.appendChild(createTopPostCard(post)))
    })
