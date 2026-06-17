import '../../core/constants/api_constants.dart';
import '../../domain/models/user.dart';
import 'api_datasource.dart';

/// Handles authentication API calls.
class AuthDataSource {
  final ApiDataSource _api;

  AuthDataSource(this._api);

  /// Register a new user.
  Future<AuthResult> register({
    required String email,
    required String password,
    required String nickname,
  }) async {
    final response = await _api.dio.post(ApiConstants.register, data: {
      'email': email,
      'password': password,
      'nickname': nickname,
    });
    return AuthResult.fromJson(response.data['data']);
  }

  /// Login with email and password.
  Future<AuthResult> login({
    required String email,
    required String password,
  }) async {
    final response = await _api.dio.post(ApiConstants.login, data: {
      'email': email,
      'password': password,
    });
    return AuthResult.fromJson(response.data['data']);
  }
}

/// Combined result of an auth operation.
class AuthResult {
  final User user;
  final String accessToken;
  final String refreshToken;

  const AuthResult({
    required this.user,
    required this.accessToken,
    required this.refreshToken,
  });

  factory AuthResult.fromJson(Map<String, dynamic> json) {
    return AuthResult(
      user: User.fromJson(json['user']),
      accessToken: json['access_token'] as String,
      refreshToken: json['refresh_token'] as String,
    );
  }
}
