import '../../core/constants/api_constants.dart';
import 'api_datasource.dart';

/// Handles daily quote API calls.
class QuoteDataSource {
  final ApiDataSource _api;

  QuoteDataSource(this._api);

  /// Get today's quote.
  Future<Map<String, dynamic>> getTodayQuote() async {
    final response = await _api.dio.get(ApiConstants.quoteToday);
    return Map<String, dynamic>.from(response.data['data'] as Map);
  }
}
