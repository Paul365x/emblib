use strict;
use warnings;
use Data::Dumper;

my @cols = (0) x 78;
my $fh;
open( $fh, "<", "jf_col.csv" ) or die "failed to open col.csv";

<$fh>;
while (my $ln = <$fh> ) {
	chomp $ln;
	my @parts = split( /\,/, $ln );
	if (scalar @parts != 5 ) {
		next;
	}
	my %col = ();
	$col{"red"} = sprintf("0x%X",$parts[0]);
	$col{"green"} = sprintf("0x%X",$parts[1]);
	$col{"blue"} = sprintf("0x%X",$parts[2]);
	$col{"name"} = $parts[3];
	$col{"name"} =~ s/\s+//g;
	$col{"index"} = $parts[4];
	#print "$parts[4]\n";
	$cols[$col{"index"}] = \%col;
}

print "var (\n";
foreach my $thr ( @cols ) {
   my $name = $thr->{ "name" };
   $name = sprintf("%s%s",$name,"Jf");
   print "\t$name\t= color.RGBA{$thr->{'red'}, $thr->{'green'}, $thr->{'blue'}, 255}\n";   
}
print ")\n\n";
print "func Janome_set() *map[string]color.Color {\n";
print "\treturn &map[string]color.Color {\n";
foreach my $thr ( @cols ) {
	my $name2 = $thr->{ "name" };
	my $name1 = sprintf("%s_%s","Jf",$name2);
	$name2 = sprintf("%s%s",$name2,"Jf");
	print "\t\t\"$name1\":\t$name2,\n";
}
print "\t}\n}\n\n";
print "func Janome_select() []color.Colour {\n";
print "\treturn []color.Color {\n";
foreach my $thr ( @cols ) {
   my $name = $thr->{ "name" };
   $name = sprintf("%s%s",$name,"Jf");
   print "\t\t$name,\n";   
}
print "\t}\n}\n";

#print Dumper(\@cols);
