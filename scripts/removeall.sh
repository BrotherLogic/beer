for c in $( cat *.txt | awk '{print substr($6,4)}'); do
    ../beer remove --id $c
done
