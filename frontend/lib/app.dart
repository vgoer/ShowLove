import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'core/constants/app_colors.dart';
import 'core/router/app_router.dart';
import 'data/datasources/api_datasource.dart';
import 'data/datasources/auth_datasource.dart';
import 'data/datasources/post_datasource.dart';
import 'data/datasources/local_storage_datasource.dart';
import 'data/repositories/auth_repository.dart';
import 'data/repositories/post_repository.dart';

/// Root widget for the Show Love app.
class ShowLoveApp extends ConsumerStatefulWidget {
  const ShowLoveApp({super.key});

  @override
  ConsumerState<ShowLoveApp> createState() => _ShowLoveAppState();
}

class _ShowLoveAppState extends ConsumerState<ShowLoveApp> {
  late final LocalStorageDataSource _storage;
  late final ApiDataSource _api;

  @override
  void initState() {
    super.initState();
    _initDependencies();
  }

  Future<void> _initDependencies() async {
    final prefs = await SharedPreferences.getInstance();
    _storage = LocalStorageDataSource(prefs);
    _api = ApiDataSource(storage: _storage);
  }

  @override
  Widget build(BuildContext context) {
    return ProviderScope(
      child: MaterialApp.router(
        title: '显出爱心',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(
            seedColor: AppColors.primary,
            surface: AppColors.background,
          ),
          scaffoldBackgroundColor: AppColors.background,
          cardColor: AppColors.card,
          appBarTheme: const AppBarTheme(
            backgroundColor: AppColors.card,
            elevation: 0,
            centerTitle: true,
          ),
          inputDecorationTheme: InputDecorationTheme(
            filled: true,
            fillColor: AppColors.card,
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: BorderSide.none,
            ),
            contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
          ),
        ),
        routerConfig: appRouter,
      ),
    );
  }
}
