const forumContainer = document.getElementById('forumContainer');
const postTpl = document.getElementById('post-template');

function getQueryParam(param) {
  const url = new URL(window.location.href);
  return url.searchParams.get(param);
}

async function loadCategory() {
  const categoryId = getQueryParam('id');
  if (!categoryId) {
    forumContainer.textContent = 'No category ID provided.';
    return;
  }

  try {
    const resp = await fetch(`http://localhost:8080/forum/api/public/feed`, { credentials: 'include' });
    if (!resp.ok) throw new Error('Failed to load data');
    const data = await resp.json();

    const category = data.categories.find(c => c.id.toString() === categoryId);
    if (!category) {
      forumContainer.textContent = 'Category not found.';
      return;
    }

    const posts = category.posts;
    forumContainer.innerHTML = `<h2>${category.name}</h2>`;

    if (posts.length === 0) {
      forumContainer.innerHTML += '<p>No posts in this category.</p>';
      return;
    }

    posts.forEach(post => {
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

      forumContainer.appendChild(postNode);
    });

  } catch (err) {
    console.error('Error:', err);
    forumContainer.textContent = 'An error occurred loading this category.';
  }
}

window.addEventListener('DOMContentLoaded', loadCategory);
