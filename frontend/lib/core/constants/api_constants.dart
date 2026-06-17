/// Backend API endpoint constants.
class ApiConstants {
  ApiConstants._();

  /// Base URL for the API Gateway.
  /// Change this to your production URL in release builds.
  static const String baseUrl = 'http://10.0.2.2:8080/api/v1'; // Android emulator → host
  // iOS simulator: 'http://localhost:8080/api/v1'
  // Real device: 'http://<your-ip>:8080/api/v1'

  // Auth
  static const String register = '/auth/register';
  static const String login = '/auth/login';
  static const String refresh = '/auth/refresh';

  // User
  static const String userMe = '/users/me';
  static const String userAvatar = '/users/me/avatar';

  // Posts
  static const String posts = '/posts';
  static String postDetail(String id) => '/posts/$id';
  static String postStickers(String id) => '/posts/$id/stickers';
  static String postReport(String id) => '/posts/$id/report';

  // Comments
  static String comments(String postId) => '/posts/$postId/comments';
  static String commentDelete(String id) => '/comments/$id';

  // Mood
  static const String moods = '/moods';
  static const String moodsWeekly = '/moods/weekly';

  // Quote
  static const String quoteToday = '/quotes/today';

  // Device
  static const String devices = '/devices';

  // Upload
  static const String uploadImage = '/upload/image';

  // Health
  static const String health = '/health';
}
