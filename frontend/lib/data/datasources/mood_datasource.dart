import '../../core/constants/api_constants.dart';
import 'api_datasource.dart';

/// Handles mood-tracking API calls.
class MoodDataSource {
  final ApiDataSource _api;

  MoodDataSource(this._api);

  /// Record today's mood.
  Future<void> recordMood({
    required int moodLevel,
    required String moodLabel,
    String note = '',
  }) async {
    await _api.dio.post(ApiConstants.moods, data: {
      'mood_level': moodLevel,
      'mood_label': moodLabel,
      'note': note,
    });
  }

  /// Get mood entries for a date range.
  Future<List<Map<String, dynamic>>> getMoods(String from, String to) async {
    final response = await _api.dio.get(ApiConstants.moods, queryParameters: {
      'from': from,
      'to': to,
    });
    final data = response.data['data'];
    return (data['entries'] as List<dynamic>)
        .map((e) => Map<String, dynamic>.from(e as Map))
        .toList();
  }

  /// Get weekly mood chart data.
  Future<List<Map<String, dynamic>>> getWeeklyMood() async {
    final response = await _api.dio.get(ApiConstants.moodsWeekly);
    final data = response.data['data'];
    return (data['points'] as List<dynamic>)
        .map((e) => Map<String, dynamic>.from(e as Map))
        .toList();
  }
}
