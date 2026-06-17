/// User model matching the backend User entity.
class User {
  final String id;
  final String email;
  final String nickname;
  final String avatarUrl;
  final String bio;
  final int kindnessScore;

  const User({
    required this.id,
    required this.email,
    required this.nickname,
    this.avatarUrl = '',
    this.bio = '',
    this.kindnessScore = 0,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] as String,
      email: json['email'] as String,
      nickname: json['nickname'] as String,
      avatarUrl: json['avatar_url'] as String? ?? '',
      bio: json['bio'] as String? ?? '',
      kindnessScore: json['kindness_score'] as int? ?? 0,
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'email': email,
        'nickname': nickname,
        'avatar_url': avatarUrl,
        'bio': bio,
        'kindness_score': kindnessScore,
      };
}
