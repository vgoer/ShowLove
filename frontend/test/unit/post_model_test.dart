import 'package:flutter_test/flutter_test.dart';
import 'package:show_love/domain/models/post.dart';
import 'package:show_love/domain/models/comment.dart';
import 'package:show_love/domain/models/user.dart';

void main() {
  group('User', () {
    test('fromJson creates valid User', () {
      final json = {
        'id': 'user-1',
        'email': 'test@example.com',
        'nickname': '小温暖',
        'avatar_url': 'http://example.com/avatar.jpg',
        'bio': '一个温暖的人',
        'kindness_score': 42,
      };

      final user = User.fromJson(json);

      expect(user.id, 'user-1');
      expect(user.email, 'test@example.com');
      expect(user.nickname, '小温暖');
      expect(user.avatarUrl, 'http://example.com/avatar.jpg');
      expect(user.kindnessScore, 42);
    });

    test('fromJson handles missing fields with defaults', () {
      final json = {
        'id': 'user-1',
        'email': 'test@example.com',
        'nickname': '小温暖',
      };

      final user = User.fromJson(json);

      expect(user.avatarUrl, '');
      expect(user.bio, '');
      expect(user.kindnessScore, 0);
    });
  });

  group('Post', () {
    test('fromJson creates valid Post', () {
      final json = {
        'id': 'post-1',
        'author_id': 'user-1',
        'author_nickname': '小温暖',
        'author_avatar': '',
        'content': '今天心情不太好',
        'mood_tag': 'sad',
        'images': ['http://img1.jpg'],
        'sticker_hug': 3,
        'sticker_cheer': 1,
        'comment_count': 5,
        'has_ai_reply': true,
        'created_at': '2026-06-17T10:00:00Z',
      };

      final post = Post.fromJson(json);

      expect(post.id, 'post-1');
      expect(post.content, '今天心情不太好');
      expect(post.moodTag, 'sad');
      expect(post.images.length, 1);
      expect(post.stickerHug, 3);
      expect(post.commentCount, 5);
      expect(post.hasAiReply, true);
    });
  });

  group('Comment', () {
    test('fromJson creates valid Comment', () {
      final json = {
        'id': 'comment-1',
        'post_id': 'post-1',
        'author_id': 'ai-bot',
        'author_nickname': '小暖',
        'content': '加油！一切都会好起来的 🌸',
        'is_ai_generated': true,
        'created_at': '2026-06-17T10:01:00Z',
      };

      final comment = Comment.fromJson(json);

      expect(comment.id, 'comment-1');
      expect(comment.isAiGenerated, true);
      expect(comment.authorNickname, '小暖');
    });
  });
}
