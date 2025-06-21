const feedURL = 'http://localhost:8080/forum/api/public/feed';

async function loadFeed() {
  try {
    const resp = await fetch(feedURL, { credentials: 'include' });
    if (!resp.ok) throw new Error('Failed to load feed');

    const data = await resp.json();

    renderFeed(data.categories || []);
  } catch (err) {
    console.error('Error loading feed:', err);
  }
}



function renderFeed(categories) {
  const container = document.getElementById('forumContainer');
  container.innerHTML = '';

  if (categories.length === 0) {
    container.textContent = 'No categories or posts available';
    return;
  }

  const categoryTpl = document.getElementById('category-template');
  const postTpl = document.getElementById('post-template');

  categories.forEach(category => {
    if (!category.posts || category.posts.length === 0) return;

    const categoryNode = categoryTpl.content.cloneNode(true);
    const titleEl = categoryNode.querySelector('.category-title');
    const link = document.createElement('a');
    link.href = `/guest/category?id=${category.id}`;
    link.textContent = category.name;
    link.classList.add('category-link'); // Optional: add a CSS class
    titleEl.appendChild(link);

    const postsContainer = categoryNode.querySelector('.category-posts');

    category.posts.forEach(post => {
      const postNode = postTpl.content.cloneNode(true);
      postNode.querySelector('.post-header').textContent = post.username || post.user_id;
      postNode.querySelector('.post-title').textContent = post.title;
      postNode.querySelector('.post-content').textContent = post.content;

      if (post.created_at) {
        postNode.querySelector('.post-time').textContent = new Date(post.created_at).toLocaleString();
      }

      const likes = post.reactions?.filter(r => r.reaction_type === 1).length || 0;
      const dislikes = post.reactions?.filter(r => r.reaction_type === 2).length || 0;

      postNode.querySelector('.like-count').textContent = likes;
      postNode.querySelector('.dislike-count').textContent = dislikes;

      postsContainer.appendChild(postNode);
    });

    container.appendChild(categoryNode);
  });
}


window.addEventListener('DOMContentLoaded', loadFeed);
