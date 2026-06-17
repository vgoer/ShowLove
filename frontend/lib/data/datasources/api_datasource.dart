import 'package:dio/dio.dart';
import '../datasources/local_storage_datasource.dart';
import '../../core/constants/api_constants.dart';

/// Dio HTTP client wrapper with JWT interceptor and automatic token refresh.
class ApiDataSource {
  late final Dio dio;
  final LocalStorageDataSource _storage;
  bool _isRefreshing = false;

  ApiDataSource({required LocalStorageDataSource storage})
      : _storage = storage {
    dio = Dio(BaseOptions(
      baseUrl: ApiConstants.baseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
      headers: {'Content-Type': 'application/json'},
    ));

    // Attach Bearer token to every request
    dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) {
        final token = _storage.getAccessToken();
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        handler.next(options);
      },
      onError: (error, handler) async {
        // Auto-refresh on 401
        if (error.response?.statusCode == 401 && !_isRefreshing) {
          _isRefreshing = true;
          final refreshed = await _tryRefresh();
          _isRefreshing = false;
          if (refreshed) {
            // Retry the original request
            final token = _storage.getAccessToken();
            error.requestOptions.headers['Authorization'] = 'Bearer $token';
            final retryResponse = await dio.fetch(error.requestOptions);
            return handler.resolve(retryResponse);
          }
        }
        handler.next(error);
      },
    ));
  }

  Future<bool> _tryRefresh() async {
    try {
      final refreshToken = _storage.getRefreshToken();
      if (refreshToken == null) return false;

      final response = await Dio(BaseOptions(
        baseUrl: ApiConstants.baseUrl,
      )).post(ApiConstants.refresh, data: {
        'refresh_token': refreshToken,
      });

      if (response.statusCode == 200) {
        final data = response.data['data'];
        await _storage.saveTokens(
          accessToken: data['access_token'],
          refreshToken: data['refresh_token'] ?? refreshToken,
        );
        return true;
      }
    } catch (_) {
      await _storage.clearTokens();
    }
    return false;
  }
}
