const feedURL = 'http://localhost:8080/forum/api/public/feed';

async function loadFeed() {
  try {
    const resp = await fetch(feedURL, { credentials: 'include' });
    if (!resp.ok) throw new Error('Failed to load feed');
    const data = await resp.json();
    renderFeed(data.posts || [], data.reactions || []);
  } catch (err) {
    console.error('Error loading feed:', err);
  }
}

function renderFeed(posts, reactions) {
  const container = document.getElementById('forumContainer');
  container.innerHTML = '';
  if (posts.length === 0) {
    container.textContent = 'No posts available';
    return;
  }
  const tpl = document.getElementById('post-template');
  posts.forEach(post => {
    const node = tpl.content.cloneNode(true);
    node.querySelector('.post-header').textContent = post.user_id;
    node.querySelector('.post-title').textContent = post.title;
    node.querySelector('.post-content').textContent = post.content;
    if (post.created_at) {
      node.querySelector('.post-time').textContent = new Date(post.created_at).toLocaleString();
    }
    const likes = reactions.filter(r => r.post_id === post.id && r.reaction_type === 1).length;
    const dislikes = reactions.filter(r => r.post_id === post.id && r.reaction_type === 2).length;
    node.querySelector('.like-count').textContent = likes;
    node.querySelector('.dislike-count').textContent = dislikes;
    container.appendChild(node);
  });
}

window.addEventListener('DOMContentLoaded', loadFeed);
