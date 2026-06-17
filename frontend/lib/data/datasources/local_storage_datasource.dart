import 'package:shared_preferences/shared_preferences.dart';

/// Manages persistent local storage for tokens and user preferences.
class LocalStorageDataSource {
  static const _accessTokenKey = 'access_token';
  static const _refreshTokenKey = 'refresh_token';
  static const _userIdKey = 'user_id';

  final SharedPreferences _prefs;

  LocalStorageDataSource(this._prefs);

  // ── Token ──

  Future<void> saveTokens({
    required String accessToken,
    required String refreshToken,
  }) async {
    await _prefs.setString(_accessTokenKey, accessToken);
    await _prefs.setString(_refreshTokenKey, refreshToken);
  }

  String? getAccessToken() => _prefs.getString(_accessTokenKey);
  String? getRefreshToken() => _prefs.getString(_refreshTokenKey);

  Future<void> clearTokens() async {
    await _prefs.remove(_accessTokenKey);
    await _prefs.remove(_refreshTokenKey);
  }

  bool get hasToken => getAccessToken() != null;

  // ── User ──

  Future<void> saveUserId(String userId) async {
    await _prefs.setString(_userIdKey, userId);
  }

  String? getUserId() => _prefs.getString(_userIdKey);

  Future<void> clearAll() async {
    await _prefs.clear();
  }
}
