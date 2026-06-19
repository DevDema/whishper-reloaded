import {
	CheckCircle2,
	Loader2,
	Clock,
	Languages,
	XCircle
} from 'lucide-svelte';

// Maps a logical status kind to its visual configuration.
// kind: 'done' | 'processing' | 'queued' | 'translating' | 'error'
export const STATUS_CONFIG = {
	done: {
		icon: CheckCircle2,
		dot: 'bg-emerald-400',
		iconColor: 'text-emerald-400',
		badge: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
		spin: false
	},
	processing: {
		icon: Loader2,
		dot: 'bg-blue-400',
		iconColor: 'text-blue-400',
		badge: 'bg-blue-500/10 text-blue-400 border-blue-500/20',
		spin: true
	},
	queued: {
		icon: Clock,
		dot: 'bg-amber-400',
		iconColor: 'text-amber-400',
		badge: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
		spin: false
	},
	translating: {
		icon: Languages,
		dot: 'bg-violet-400',
		iconColor: 'text-violet-400',
		badge: 'bg-violet-500/10 text-violet-400 border-violet-500/20',
		spin: false
	},
	error: {
		icon: XCircle,
		dot: 'bg-red-400',
		iconColor: 'text-red-400',
		badge: 'bg-red-500/10 text-red-400 border-red-500/20',
		spin: false
	}
};
