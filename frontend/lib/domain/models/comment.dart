/// Comment model matching the backend Comment entity.
class Comment {
  final String id;
  final String postId;
  final String authorId;
  final String authorNickname;
  final String authorAvatar;
  final String content;
  final bool isAiGenerated;
  final DateTime createdAt;

  const Comment({
    required this.id,
    required this.postId,
    required this.authorId,
    required this.authorNickname,
    this.authorAvatar = '',
    required this.content,
    this.isAiGenerated = false,
    required this.createdAt,
  });

  factory Comment.fromJson(Map<String, dynamic> json) {
    return Comment(
      id: json['id'] as String,
      postId: json['post_id'] as String,
      authorId: json['author_id'] as String,
      authorNickname: json['author_nickname'] as String,
      authorAvatar: json['author_avatar'] as String? ?? '',
      content: json['content'] as String,
      isAiGenerated: json['is_ai_generated'] as bool? ?? false,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }
}
