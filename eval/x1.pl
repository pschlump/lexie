#!/usr/bin/perl

$n = 1;
while ( <> ) {
	if ( /"Error \([^)]*\):/ ) {
		$sn = sprintf ( "%05d", $n );
		s/Error \([^)]*\):/Error (Eval$sn):/;
		$n++;
	}
	print $_;
}
