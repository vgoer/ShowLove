import '../../core/constants/api_constants.dart';
import '../../domain/models/post.dart';
import 'api_datasource.dart';

/// Handles post-related API calls.
class PostDataSource {
  final ApiDataSource _api;

  PostDataSource(this._api);

  /// Fetch paginated posts.
  Future<PostListResult> listPosts({
    String sort = 'latest',
    int page = 1,
    int pageSize = 20,
  }) async {
    final response = await _api.dio.get(ApiConstants.posts, queryParameters: {
      'sort': sort,
      'page': page,
      'size': pageSize,
    });
    final data = response.data['data'];
    return PostListResult(
      posts: (data['posts'] as List<dynamic>)
          .map((e) => Post.fromJson(e as Map<String, dynamic>))
          .toList(),
      total: data['total'] as int,
      page: data['page'] as int,
    );
  }

  /// Create a new post.
  Future<Post> createPost({
    required String content,
    required String moodTag,
    List<String> images = const [],
  }) async {
    final response = await _api.dio.post(ApiConstants.posts, data: {
      'content': content,
      'mood_tag': moodTag,
      'images': images,
    });
    return Post.fromJson(response.data['data']);
  }

  /// Get a single post by ID.
  Future<Post> getPost(String postId) async {
    final response = await _api.dio.get(ApiConstants.postDetail(postId));
    return Post.fromJson(response.data['data']);
  }

  /// Send a sticker to a post.
  Future<Post> sendSticker(String postId, String stickerType) async {
    final response =
        await _api.dio.post(ApiConstants.postStickers(postId), data: {
      'sticker_type': stickerType,
    });
    return Post.fromJson(response.data['data']);
  }

  /// Report a post.
  Future<void> reportPost(String postId, String reason) async {
    await _api.dio.post(ApiConstants.postReport(postId), data: {
      'reason': reason,
    });
  }

  /// Delete own post.
  Future<void> deletePost(String postId) async {
    await _api.dio.delete(ApiConstants.postDetail(postId));
  }
}

class PostListResult {
  final List<Post> posts;
  final int total;
  final int page;

  const PostListResult({
    required this.posts,
    required this.total,
    required this.page,
  });
}
