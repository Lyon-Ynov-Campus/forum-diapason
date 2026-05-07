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
    authorEl.textContent = post.author
    authorEl.href = `/profile?id=${post.author_id}`
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
    initPostCard(article, post)
    return card
}

const profileId = parseInt(new URLSearchParams(window.location.search).get('id')) || 1

Promise.all([
    fetch('/data/profile.json').then(r => r.json()),
    fetch('/data/posts.json').then(r => r.json())
]).then(([profileData, posts]) => {
    const profiles = Array.isArray(profileData) ? profileData : [profileData]
    const profile = profiles.find(p => p.id === profileId) || profiles[0]
    if (!profile) return

    document.getElementById('profile-pseudo').textContent = profile.pseudo
    document.querySelector('#profile-ville span:last-child').textContent = profile.ville
    document.getElementById('profile-followers').textContent = profile.followers
    document.getElementById('profile-posts').textContent = profile.posts

    const userPosts = posts.filter(p => p.author_id === profile.id)
    const container = document.getElementById('profile-posts-container')
    userPosts.forEach(post => container.appendChild(createPostCard(post)))

    // Bouton edit → ouvre le modal pré-rempli
    document.getElementById('open-edit-profile-btn')?.addEventListener('click', () => {
        openEditProfileModal(profile)
    })
})
