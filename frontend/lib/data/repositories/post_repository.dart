import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../datasources/post_datasource.dart';
import '../../domain/models/post.dart';

/// Manages post data operations.
class PostRepository {
  final PostDataSource _remote;

  PostRepository(this._remote);

  Future<PostListResult> getPosts({
    String sort = 'latest',
    int page = 1,
    int pageSize = 20,
  }) async {
    return _remote.listPosts(sort: sort, page: page, pageSize: pageSize);
  }

  Future<Post> createPost({
    required String content,
    required String moodTag,
    List<String> images = const [],
  }) async {
    return _remote.createPost(
      content: content,
      moodTag: moodTag,
      images: images,
    );
  }

  Future<Post> getPost(String postId) => _remote.getPost(postId);

  Future<Post> sendSticker(String postId, String type) =>
      _remote.sendSticker(postId, type);

  Future<void> reportPost(String postId, String reason) =>
      _remote.reportPost(postId, reason);

  Future<void> deletePost(String postId) => _remote.deletePost(postId);
}

final postRepositoryProvider = Provider<PostRepository>((ref) {
  throw UnimplementedError('PostRepository must be provided');
});
