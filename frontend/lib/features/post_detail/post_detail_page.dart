import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/constants/app_colors.dart';

/// Post detail page showing full post content, comments, and sticker actions.
class PostDetailPage extends ConsumerStatefulWidget {
  final String postId;

  const PostDetailPage({required this.postId, super.key});

  @override
  ConsumerState<PostDetailPage> createState() => _PostDetailPageState();
}

class _PostDetailPageState extends ConsumerState<PostDetailPage> {
  final _commentCtrl = TextEditingController();
  bool _loading = false;

  @override
  void dispose() {
    _commentCtrl.dispose();
    super.dispose();
  }

  void _sendSticker(String type) {
    // TODO: Call PostRepository.sendSticker
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('贴纸已发送 🌸'), duration: Duration(seconds: 1)),
    );
  }

  Future<void> _sendComment() async {
    if (_commentCtrl.text.trim().isEmpty) return;
    setState(() => _loading = true);
    try {
      // TODO: Call CommentDataSource.createComment
      _commentCtrl.clear();
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('帖子详情'),
        backgroundColor: Colors.transparent,
        elevation: 0,
        foregroundColor: AppColors.textPrimary,
      ),
      body: Column(
        children: [
          // Post content area (placeholder)
          Expanded(
            child: ListView(
              padding: const EdgeInsets.all(16),
              children: const [
                // Post content card
                Card(
                  color: AppColors.card,
                  child: Padding(
                    padding: EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text('加载中...', style: TextStyle(color: AppColors.textSecondary)),
                      ],
                    ),
                  ),
                ),
                SizedBox(height: 16),
                // AI reply banner
                Card(
                  color: AppColors.lavender,
                  child: Padding(
                    padding: EdgeInsets.all(12),
                    child: Row(
                      children: [
                        Icon(Icons.auto_awesome, color: AppColors.primary, size: 20),
                        SizedBox(width: 8),
                        Expanded(
                          child: Text('小暖正在为你生成暖心回复...',
                            style: TextStyle(color: AppColors.textSecondary, fontSize: 13)),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ),
          // Sticker bar
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            decoration: const BoxDecoration(
              color: AppColors.card,
              border: Border(top: BorderSide(color: AppColors.divider)),
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: [
                _StickerButton(emoji: '🤗', label: '抱抱', onTap: () => _sendSticker('hug')),
                _StickerButton(emoji: '💪', label: '加油', onTap: () => _sendSticker('cheer')),
                _StickerButton(emoji: '💛', label: '我懂你', onTap: () => _sendSticker('understand')),
              ],
            ),
          ),
          // Comment input
          Container(
            padding: const EdgeInsets.all(12),
            decoration: const BoxDecoration(
              color: AppColors.card,
              border: Border(top: BorderSide(color: AppColors.divider)),
            ),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _commentCtrl,
                    decoration: const InputDecoration(
                      hintText: '写下温暖的评论...',
                      border: InputBorder.none,
                    ),
                  ),
                ),
                IconButton(
                  icon: const Icon(Icons.send, color: AppColors.primary),
                  onPressed: _loading ? null : _sendComment,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _StickerButton extends StatelessWidget {
  final String emoji;
  final String label;
  final VoidCallback onTap;

  const _StickerButton({
    required this.emoji,
    required this.label,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(emoji, style: const TextStyle(fontSize: 28)),
          const SizedBox(height: 4),
          Text(label, style: const TextStyle(fontSize: 11, color: AppColors.textSecondary)),
        ],
      ),
    );
  }
}
