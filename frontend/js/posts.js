function timeAgo(dateStr) {
    const normalized = dateStr.replace(' +0000 UTC', 'Z').replace(' ', 'T')
    const diff = Math.floor((Date.now() - new Date(normalized)) / 60000)
    if (diff < 1) return "à l'instant"
    if (diff < 60) return `il y a ${diff} minute${diff > 1 ? 's' : ''}`
    const h = Math.floor(diff / 60)
    if (h < 24) return `il y a ${h} heure${h > 1 ? 's' : ''}`
    const d = Math.floor(h / 24)
    return `il y a ${d} jour${d > 1 ? 's' : ''}`
}

function createPostCard(post) {
    const tpl = document.getElementById('post-card').content.cloneNode(true)
    const authorEl = tpl.querySelector('.post-author')
    authorEl.textContent = post.author_pseudo
    authorEl.href = `/profile?id=${post.user_id}`
    tpl.querySelector('.post-title').textContent = post.titre
    tpl.querySelector('.post-content').textContent = post.contenu
    tpl.querySelector('.post-date').textContent = timeAgo(post.date_publication)
    tpl.querySelector('.post-tags').textContent = (post.tags || []).map(t => `#${t}`).join(' ')
    tpl.querySelector('.post-likes').textContent = post.like_count

    // Avatar auteur
    const photoImg = tpl.querySelector('.post-author-photo')
    const photoDefault = tpl.querySelector('.post-author-avatar-default')
    if (post.author_photo && post.author_photo !== '' && photoImg) {
        photoImg.src = post.author_photo.startsWith('http')
            ? post.author_photo
            : `http://localhost:8080${post.author_photo}`
        photoImg.classList.remove('hidden')
        if (photoDefault) photoDefault.classList.add('hidden')
    }

    const article = tpl.querySelector('article')
    article.style.cursor = 'pointer'
    article.addEventListener('click', (e) => {
        if (e.target.closest('button')) return
        window.location.href = `/post?id=${post.id}`
    })
    initPostCard(article, post)
    return tpl
}

function createTopPostCard(post) {
    const tpl = document.getElementById('top-post-card').content.cloneNode(true)
    tpl.querySelector('.top-post-author').textContent = post.author_pseudo
    tpl.querySelector('.top-post-content').textContent = post.contenu
    tpl.querySelector('.top-post-date').textContent = timeAgo(post.date_publication)
    const article = tpl.querySelector('article')
    article.style.cursor = 'pointer'
    article.addEventListener('click', () => {
        window.location.href = `/post?id=${post.id}`
    })
    return tpl
}

function loadPosts(tagFilter = '') {
    const container = document.getElementById('posts-container')
    container.innerHTML = ''

    const url = tagFilter
        ? `${API}/api/posts/tag/${encodeURIComponent(tagFilter)}`
        : `${API}/api/posts`

    fetch(url, { credentials: 'include' })
        .then(r => r.json())
        .then(posts => {
            if (!posts || posts.length === 0) {
                container.innerHTML = '<p class="text-sm text-gray-400 text-center py-8">Aucun post trouvé.</p>'
                return
            }
            posts.forEach(post => container.appendChild(createPostCard(post)))
        })
        .catch(() => {
            container.innerHTML = '<p class="text-sm text-red-400 text-center py-8">Erreur de chargement.</p>'
        })
}

function reloadPosts() {
    const searchInput = document.querySelector('header input[type="text"]')
    const currentTag = searchInput ? searchInput.value.trim().replace(/^#/, '') : ''
    loadPosts(currentTag)
}

//loadPosts()

function loadTopPosts() {
    fetch(`${API}/api/posts/top?limit=6`, { credentials: 'include' })
        .then(r => r.json())
        .then(posts => {
            const container = document.getElementById('top-posts-container')
            if (!container || !posts) return
            posts.forEach(post => container.appendChild(createTopPostCard(post)))
        })
}

function initSearch() {
    const input = document.querySelector('header input[type="text"]')
    
    if (!input) return

    input.addEventListener('keydown', (e) => {
        if (e.key !== 'Enter') return
        const val = input.value.trim().replace(/^#/, '')
        loadPosts(val)
        const url = val ? `/?tag=${encodeURIComponent(val)}` : '/'
        window.history.pushState({}, '', url)
    })

    const allHeaderBtns = document.querySelectorAll('header button')
    allHeaderBtns.forEach(btn => {
        const paths = btn.querySelectorAll('svg path')
        const hasSearch = Array.from(paths).some(p => p.getAttribute('d')?.includes('M21 21'))
        if (hasSearch) {
            btn.addEventListener('click', () => {
                const val = input.value.trim().replace(/^#/, '')
                loadPosts(val)
            })
        }
    })                          

    const urlParams = new URLSearchParams(window.location.search)
    const tagFromURL = urlParams.get('tag')
    if (tagFromURL) {
        input.value = `#${tagFromURL}`
        loadPosts(tagFromURL)
        return
    }
}

document.addEventListener('DOMContentLoaded', () => {
    loadPosts()
    loadTopPosts()
    initSearch()
})