package config

import ()

const IdxPostsUserID = `CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);`
const IdxPostCategoriesPostID = `CREATE INDEX IF NOT EXISTS idx_post_categories_post_id ON post_categories(post_id);`
const IdxPostCategoriesCategoryID = `CREATE INDEX IF NOT EXISTS idx_post_categories_category_id ON post_categories(category_id);`
const IdxCommentsPostID = `CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);`
const IdxCommentsUserID = `CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);`
const IdxReactionsUserID = `CREATE INDEX IF NOT EXISTS idx_reactions_user_id ON reactions(user_id);`
const IdxReactionsPostID = `CREATE INDEX IF NOT EXISTS idx_reactions_post_id ON reactions(post_id);`
const IdxReactionsCommentID = `CREATE INDEX IF NOT EXISTS idx_reactions_comment_id ON reactions(comment_id);`