function timeAgo(dateStr) {
    const diff = Math.floor((Date.now() - new Date(dateStr)) / 60000)
    if (diff < 60) return `il y a ${diff} minute${diff > 1 ? 's' : ''}`
    const h = Math.floor(diff / 60)
    if (h < 24) return `il y a ${h} heure${h > 1 ? 's' : ''}`
    const d = Math.floor(h / 24)
    return `il y a ${d} jour${d > 1 ? 's' : ''}`
}

function profileUrl(pseudo) {
    return `/profile?pseudo=${encodeURIComponent(pseudo)}`
}

function createReply(reply) {
    const tpl = document.getElementById('reply-card').content.cloneNode(true)
    const a = tpl.querySelector('.reply-author')
    a.textContent = reply.author
    a.href = profileUrl(reply.author)
    tpl.querySelector('.reply-content').textContent = reply.content
    tpl.querySelector('.reply-date').textContent = timeAgo(reply.created_at)
    return tpl
}

function createComment(comment) {
    const tpl = document.getElementById('comment-card').content.cloneNode(true)
    const authorA = tpl.querySelector('.comment-author')
    authorA.textContent = comment.author
    authorA.href = profileUrl(comment.author)
    tpl.querySelector('.comment-content').textContent = comment.content
    tpl.querySelector('.comment-date').textContent = timeAgo(comment.created_at)

    // Like commentaire
    let liked = false
    let count = comment.likes || 0
    const likesEl  = tpl.querySelector('.comment-likes')
    const likeBtn  = tpl.querySelector('.comment-like-btn')
    likesEl.textContent = count
    likeBtn.addEventListener('click', () => {
        liked = !liked
        count += liked ? 1 : -1
        likesEl.textContent = count
        likeBtn.classList.toggle('text-red-500', liked)
        const svg = likeBtn.querySelector('svg')
        svg.setAttribute('fill', liked ? '#ef4444' : 'none')
        svg.setAttribute('stroke', liked ? '#ef4444' : 'currentColor')
    })

    // Réponses
    const repliesEl = tpl.querySelector('.comment-replies')
    ;(comment.replies || []).forEach(r => repliesEl.appendChild(createReply(r)))

    // Toggle formulaire réponse
    const replyBtn  = tpl.querySelector('.comment-reply-btn')
    const replyForm = tpl.querySelector('.comment-reply-form')
    replyBtn.addEventListener('click', () => replyForm.classList.toggle('hidden'))

    tpl.querySelector('.comment-reply-submit').addEventListener('click', () => {
        const input = replyForm.querySelector('input')
        const val = input.value.trim()
        if (!val) return
        repliesEl.appendChild(createReply({ author: 'moi', content: val, created_at: new Date().toISOString() }))
        input.value = ''
        replyForm.classList.add('hidden')
    })

    return tpl
}

const postId = parseInt(new URLSearchParams(window.location.search).get('id'))

Promise.all([
    fetch('/data/posts.json').then(r => r.json()),
    fetch('/data/post-detail.json').then(r => r.json())
]).then(([posts, detail]) => {
    const post = postId ? (posts.find(p => p.id === postId) || detail) : detail
    if (post.id === detail.id) post.comments = detail.comments
    else post.comments = []
    render(post)
})

function render(post) {
    const postAuthorEl = document.getElementById('post-author')
    postAuthorEl.textContent = post.author
    postAuthorEl.href = profileUrl(post.author)
    document.getElementById('post-title').textContent   = post.title
    document.getElementById('post-content').textContent = post.content
    document.getElementById('post-date').textContent    = timeAgo(post.created_at)
    document.getElementById('post-tags').textContent    = (post.tags || []).map(t => `#${t}`).join(' ')

    // Like du post
    let liked = false
    let count = post.likes
    const likesEl = document.getElementById('post-likes')
    const likeBtn = document.querySelector('.post-detail-like')
    const heart   = document.querySelector('.post-detail-heart')
    likesEl.textContent = count
    likeBtn?.addEventListener('click', () => {
        liked = !liked
        count += liked ? 1 : -1
        likesEl.textContent = count
        heart.setAttribute('fill', liked ? '#ef4444' : 'none')
        heart.setAttribute('stroke', liked ? '#ef4444' : 'currentColor')
        likeBtn.classList.toggle('text-red-500', liked)
    })

    // Share du post
    document.querySelector('.post-detail-share')?.addEventListener('click', () => {
        navigator.clipboard.writeText(location.href)
            .then(() => showToast('Lien copié !'))
    })

    // Commentaires
    const list = document.getElementById('comments-list')
    post.comments.forEach(c => list.appendChild(createComment(c)))

    // Nouveau commentaire
    document.getElementById('new-comment-submit').addEventListener('click', () => {
        const input = document.getElementById('new-comment-input')
        const val = input.value.trim()
        if (!val) return
        list.appendChild(createComment({
            id: Date.now(), author: 'moi', content: val,
            created_at: new Date().toISOString(), likes: 0, replies: []
        }))
        input.value = ''
        list.scrollTop = list.scrollHeight
    })

    // Enter pour soumettre
    document.getElementById('new-comment-input')?.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') document.getElementById('new-comment-submit').click()
    })
}
