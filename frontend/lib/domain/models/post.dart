/// Post model matching the backend Post entity.
class Post {
  final String id;
  final String authorId;
  final String authorNickname;
  final String authorAvatar;
  final String content;
  final String moodTag;
  final List<String> images;
  final String? voiceUrl;
  final int stickerHug;
  final int stickerCheer;
  final int stickerUnderstand;
  final int commentCount;
  final bool hasAiReply;
  final DateTime createdAt;

  const Post({
    required this.id,
    required this.authorId,
    required this.authorNickname,
    this.authorAvatar = '',
    required this.content,
    required this.moodTag,
    this.images = const [],
    this.voiceUrl,
    this.stickerHug = 0,
    this.stickerCheer = 0,
    this.stickerUnderstand = 0,
    this.commentCount = 0,
    this.hasAiReply = false,
    required this.createdAt,
  });

  factory Post.fromJson(Map<String, dynamic> json) {
    return Post(
      id: json['id'] as String,
      authorId: json['author_id'] as String,
      authorNickname: json['author_nickname'] as String,
      authorAvatar: json['author_avatar'] as String? ?? '',
      content: json['content'] as String,
      moodTag: json['mood_tag'] as String,
      images: (json['images'] as List<dynamic>?)
              ?.map((e) => e.toString())
              .toList() ??
          [],
      voiceUrl: json['voice_url'] as String?,
      stickerHug: json['sticker_hug'] as int? ?? 0,
      stickerCheer: json['sticker_cheer'] as int? ?? 0,
      stickerUnderstand: json['sticker_understand'] as int? ?? 0,
      commentCount: json['comment_count'] as int? ?? 0,
      hasAiReply: json['has_ai_reply'] as bool? ?? false,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }
}
