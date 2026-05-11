function timeAgo(dateStr) {
    const diff = Math.floor((Date.now() - new Date(dateStr)) / 60000)
    if (diff < 60) return `il y a ${diff} minute${diff > 1 ? 's' : ''}`
    const h = Math.floor(diff / 60)
    if (h < 24) return `il y a ${h} heure${h > 1 ? 's' : ''}`
    const d = Math.floor(h / 24)
    return `il y a ${d} jour${d > 1 ? 's' : ''}`
}

function createReply(reply) {
    const tpl = document.getElementById('reply-card').content.cloneNode(true)
    tpl.querySelector('.reply-author').textContent = reply.author
    tpl.querySelector('.reply-content').textContent = reply.content
    tpl.querySelector('.reply-date').textContent = timeAgo(reply.created_at)
    return tpl
}

function createComment(comment) {
    const tpl = document.getElementById('comment-card').content.cloneNode(true)
    tpl.querySelector('.comment-author').textContent = comment.author
    tpl.querySelector('.comment-content').textContent = comment.content
    tpl.querySelector('.comment-date').textContent = timeAgo(comment.created_at)
    tpl.querySelector('.comment-likes').textContent = comment.likes || 0

    const repliesEl = tpl.querySelector('.comment-replies')
    ;(comment.replies || []).forEach(r => repliesEl.appendChild(createReply(r)))

    const replyBtn  = tpl.querySelector('.comment-reply-btn')
    const replyForm = tpl.querySelector('.comment-reply-form')
    replyBtn.addEventListener('click', () => replyForm.classList.toggle('hidden'))

    tpl.querySelector('.comment-reply-submit').addEventListener('click', () => {
        const input = replyForm.querySelector('input')
        const val = input.value.trim()
        if (!val) return
        const fakeReply = { author: 'moi', content: val, created_at: new Date().toISOString() }
        repliesEl.appendChild(createReply(fakeReply))
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
    // Fusionne les commentaires du detail si même post
    if (post.id === detail.id) post.comments = detail.comments
    else post.comments = []
    render(post)
})

function render(post) {
        document.getElementById('post-author').textContent  = post.author
        document.getElementById('post-title').textContent   = post.title
        document.getElementById('post-content').textContent = post.content
        document.getElementById('post-date').textContent    = timeAgo(post.created_at)
        document.getElementById('post-tags').textContent    = (post.tags || []).map(t => `#${t}`).join(' ')
        document.getElementById('post-likes').textContent   = post.likes

        const list = document.getElementById('comments-list')
        post.comments.forEach(c => list.appendChild(createComment(c)))

        document.getElementById('new-comment-submit').addEventListener('click', () => {
            const input = document.getElementById('new-comment-input')
            const val = input.value.trim()
            if (!val) return
            const fakeComment = {
                id: Date.now(), author: 'moi', content: val,
                created_at: new Date().toISOString(), likes: 0, replies: []
            }
            list.appendChild(createComment(fakeComment))
            input.value = ''
            list.scrollTop = list.scrollHeight
        })
}
