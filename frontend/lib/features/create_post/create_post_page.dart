import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../core/constants/app_colors.dart';

/// Page for creating a new post with content, mood tag, and optional images.
class CreatePostPage extends ConsumerStatefulWidget {
  const CreatePostPage({super.key});

  @override
  ConsumerState<CreatePostPage> createState() => _CreatePostPageState();
}

class _CreatePostPageState extends ConsumerState<CreatePostPage> {
  final _contentCtrl = TextEditingController();
  String _selectedMood = 'sad';
  bool _submitting = false;

  static const _moods = {
    'sad': '😢 难过',
    'anxious': '😰 焦虑',
    'lonely': '🥺 孤独',
    'stressed': '😫 压力',
    'angry': '😤 烦躁',
    'confused': '😕 迷茫',
    'calm': '😊 平静',
    'happy': '😄 开心',
  };

  @override
  void dispose() {
    _contentCtrl.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (_contentCtrl.text.trim().isEmpty) return;
    setState(() => _submitting = true);
    try {
      // TODO: Call PostRepository.createPost
      if (mounted) context.pop();
    } finally {
      if (mounted) setState(() => _submitting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('倾诉心事'),
        backgroundColor: Colors.transparent,
        elevation: 0,
        foregroundColor: AppColors.textPrimary,
        actions: [
          TextButton(
            onPressed: _submitting ? null : _submit,
            child: const Text('发布', style: TextStyle(
              color: AppColors.primary, fontWeight: FontWeight.bold, fontSize: 16)),
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Mood tag selector
            const Text('此刻心情', style: TextStyle(
              color: AppColors.textSecondary, fontSize: 14)),
            const SizedBox(height: 8),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: _moods.entries.map((e) {
                final selected = _selectedMood == e.key;
                return GestureDetector(
                  onTap: () => setState(() => _selectedMood = e.key),
                  child: Container(
                    padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
                    decoration: BoxDecoration(
                      color: selected ? AppColors.primary.withAlpha(30) : AppColors.card,
                      borderRadius: BorderRadius.circular(20),
                      border: Border.all(
                        color: selected ? AppColors.primary : AppColors.divider,
                      ),
                    ),
                    child: Text(e.value, style: TextStyle(
                      color: selected ? AppColors.primary : AppColors.textSecondary,
                    )),
                  ),
                );
              }).toList(),
            ),
            const SizedBox(height: 24),
            // Content input
            TextField(
              controller: _contentCtrl,
              maxLines: 8,
              maxLength: 5000,
              decoration: const InputDecoration(
                hintText: '写下你想说的话...',
                hintStyle: TextStyle(color: AppColors.textSecondary),
                border: OutlineInputBorder(
                  borderSide: BorderSide.none,
                ),
                filled: true,
                fillColor: AppColors.card,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
