import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../core/constants/app_colors.dart';

/// User profile page showing info, stats, and settings.
class ProfilePage extends ConsumerWidget {
  const ProfilePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('我的'),
        backgroundColor: Colors.transparent,
        elevation: 0,
        foregroundColor: AppColors.textPrimary,
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // Avatar & name
          const Center(
            child: Column(
              children: [
                CircleAvatar(
                  radius: 50,
                  backgroundColor: AppColors.primary,
                  child: Icon(Icons.person, size: 50, color: Colors.white),
                ),
                SizedBox(height: 12),
                Text('小温暖', style: TextStyle(
                  fontSize: 20, fontWeight: FontWeight.bold, color: AppColors.textPrimary)),
                SizedBox(height: 4),
                Text('善意积分: 0', style: TextStyle(color: AppColors.textSecondary)),
              ],
            ),
          ),
          const SizedBox(height: 32),
          // Menu items
          _MenuTile(
            icon: Icons.edit,
            title: '编辑资料',
            onTap: () {},
          ),
          _MenuTile(
            icon: Icons.favorite,
            title: '每日语录',
            onTap: () {},
          ),
          _MenuTile(
            icon: Icons.settings,
            title: '设置',
            onTap: () {},
          ),
          const Divider(height: 32, color: AppColors.divider),
          _MenuTile(
            icon: Icons.logout,
            title: '退出登录',
            textColor: AppColors.error,
            onTap: () {
              // TODO: AuthRepository.logout
              context.go('/login');
            },
          ),
        ],
      ),
    );
  }
}

class _MenuTile extends StatelessWidget {
  final IconData icon;
  final String title;
  final VoidCallback onTap;
  final Color? textColor;

  const _MenuTile({
    required this.icon,
    required this.title,
    required this.onTap,
    this.textColor,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Icon(icon, color: textColor ?? AppColors.textSecondary),
      title: Text(title, style: TextStyle(color: textColor ?? AppColors.textPrimary)),
      trailing: const Icon(Icons.chevron_right, color: AppColors.divider),
      onTap: onTap,
    );
  }
}
