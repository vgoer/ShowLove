import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../datasources/auth_datasource.dart';
import '../datasources/local_storage_datasource.dart';
import '../../domain/models/user.dart';

/// Manages authentication state and operations.
class AuthRepository {
  final AuthDataSource _remote;
  final LocalStorageDataSource _local;

  AuthRepository(this._remote, this._local);

  bool get isLoggedIn => _local.hasToken;

  String? get userId => _local.getUserId();

  /// Register and auto-login.
  Future<User> register({
    required String email,
    required String password,
    required String nickname,
  }) async {
    final result = await _remote.register(
      email: email,
      password: password,
      nickname: nickname,
    );
    await _persistAuth(result.user, result.accessToken, result.refreshToken);
    return result.user;
  }

  /// Login with credentials.
  Future<User> login({
    required String email,
    required String password,
  }) async {
    final result = await _remote.login(email: email, password: password);
    await _persistAuth(result.user, result.accessToken, result.refreshToken);
    return result.user;
  }

  /// Logout and clear stored tokens.
  Future<void> logout() async {
    await _local.clearAll();
  }

  Future<void> _persistAuth(
    User user,
    String accessToken,
    String refreshToken,
  ) async {
    await _local.saveTokens(
      accessToken: accessToken,
      refreshToken: refreshToken,
    );
    await _local.saveUserId(user.id);
  }
}

/// Riverpod provider for AuthRepository.
final authRepositoryProvider = Provider<AuthRepository>((ref) {
  // Will be overridden in app.dart with actual instances
  throw UnimplementedError('AuthRepository must be provided');
});
