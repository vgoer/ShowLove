import '../../core/constants/api_constants.dart';
import '../../domain/models/comment.dart';
import 'api_datasource.dart';

/// Handles comment-related API calls.
class CommentDataSource {
  final ApiDataSource _api;

  CommentDataSource(this._api);

  /// List comments for a post.
  Future<List<Comment>> listComments(String postId,
      {int page = 1, int pageSize = 20}) async {
    final response = await _api.dio.get(
      ApiConstants.comments(postId),
      queryParameters: {'page': page, 'size': pageSize},
    );
    final data = response.data['data'];
    return (data['comments'] as List<dynamic>)
        .map((e) => Comment.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  /// Create a comment on a post.
  Future<Comment> createComment(String postId, String content) async {
    final response =
        await _api.dio.post(ApiConstants.comments(postId), data: {
      'content': content,
    });
    return Comment.fromJson(response.data['data']);
  }

  /// Delete own comment.
  Future<void> deleteComment(String commentId) async {
    await _api.dio.delete(ApiConstants.commentDelete(commentId));
  }
}
