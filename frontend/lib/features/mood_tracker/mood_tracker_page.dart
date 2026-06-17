import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/constants/app_colors.dart';

/// Mood tracker page with 1-10 thermometer and weekly chart.
class MoodTrackerPage extends ConsumerStatefulWidget {
  const MoodTrackerPage({super.key});

  @override
  ConsumerState<MoodTrackerPage> createState() => _MoodTrackerPageState();
}

class _MoodTrackerPageState extends ConsumerState<MoodTrackerPage> {
  int _moodLevel = 5;
  String _moodLabel = '平静';

  static const _labels = ['难过', '焦虑', '烦躁', '迷茫', '平静', '还行', '不错', '开心', '很棒', '超赞'];

  void _saveMood() {
    // TODO: Call MoodDataSource.recordMood
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('心情已记录 🌸'), duration: Duration(seconds: 1)),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('情绪温度计'),
        backgroundColor: Colors.transparent,
        elevation: 0,
        foregroundColor: AppColors.textPrimary,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Column(
          children: [
            const SizedBox(height: 40),
            // Mood level display
            Text(_moodLevel.toString(), style: TextStyle(
              fontSize: 72, fontWeight: FontWeight.bold,
              color: _moodColor(_moodLevel),
            )),
            Text(_labels[_moodLevel - 1], style: const TextStyle(
              fontSize: 20, color: AppColors.textSecondary)),
            const SizedBox(height: 32),
            // Slider
            SliderTheme(
              data: SliderThemeData(
                activeTrackColor: _moodColor(_moodLevel),
                thumbColor: _moodColor(_moodLevel),
              ),
              child: Slider(
                value: _moodLevel.toDouble(),
                min: 1, max: 10, divisions: 9,
                onChanged: (v) => setState(() {
                  _moodLevel = v.round();
                  _moodLabel = _labels[_moodLevel - 1];
                }),
              ),
            ),
            const SizedBox(height: 40),
            // Save button
            SizedBox(
              width: double.infinity,
              height: 48,
              child: ElevatedButton(
                onPressed: _saveMood,
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppColors.primary,
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                ),
                child: const Text('记录今日心情', style: TextStyle(fontSize: 16)),
              ),
            ),
            const SizedBox(height: 40),
            // Weekly chart placeholder
            Container(
              height: 200,
              decoration: BoxDecoration(
                color: AppColors.card,
                borderRadius: BorderRadius.circular(16),
              ),
              child: const Center(
                child: Text('本周情绪曲线',
                    style: TextStyle(color: AppColors.textSecondary)),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Color _moodColor(int level) {
    if (level <= 3) return const Color(0xFF6B9BD2);   // blue
    if (level <= 5) return const Color(0xFF98D8C8);   // mint
    if (level <= 7) return const Color(0xFFFFE4B5);  // cream
    return AppColors.primary;                          // coral
  }
}
