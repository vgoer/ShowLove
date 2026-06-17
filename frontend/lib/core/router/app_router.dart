import 'package:go_router/go_router.dart';
import '../../features/auth/login_page.dart';
import '../../features/auth/register_page.dart';
import '../../features/feed/feed_page.dart';
import '../../features/post_detail/post_detail_page.dart';
import '../../features/create_post/create_post_page.dart';
import '../../features/mood_tracker/mood_tracker_page.dart';
import '../../features/profile/profile_page.dart';

/// Application route configuration using go_router.
final appRouter = GoRouter(
  initialLocation: '/feed',
  routes: [
    GoRoute(
      path: '/login',
      name: 'login',
      builder: (context, state) => const LoginPage(),
    ),
    GoRoute(
      path: '/register',
      name: 'register',
      builder: (context, state) => const RegisterPage(),
    ),
    GoRoute(
      path: '/feed',
      name: 'feed',
      builder: (context, state) => const FeedPage(),
    ),
    GoRoute(
      path: '/post/create',
      name: 'create-post',
      builder: (context, state) => const CreatePostPage(),
    ),
    GoRoute(
      path: '/post/:id',
      name: 'post-detail',
      builder: (context, state) =>
          PostDetailPage(postId: state.pathParameters['id']!),
    ),
    GoRoute(
      path: '/mood',
      name: 'mood-tracker',
      builder: (context, state) => const MoodTrackerPage(),
    ),
    GoRoute(
      path: '/profile',
      name: 'profile',
      builder: (context, state) => const ProfilePage(),
    ),
  ],
);
