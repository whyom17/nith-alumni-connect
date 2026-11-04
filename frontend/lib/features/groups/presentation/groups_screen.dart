import 'package:flutter/material.dart';
import '../../../../core/widgets/coming_soon.dart';

/// The Groups tab screen, rendered as a Coming Soon placeholder for Phase 1.
class GroupsScreen extends StatelessWidget {
  const GroupsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Groups & Communities'),
      ),
      body: const ComingSoon(
        title: 'Communities Coming Soon',
        message:
            'Interest-based groups, chapters, and discussions will launch in Phase 2.',
        icon: Icons.groups_2_rounded,
      ),
    );
  }
}
